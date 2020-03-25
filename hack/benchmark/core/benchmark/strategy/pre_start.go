package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

type PreStart func(context.Context, *testing.B, interface{}, assets.Dataset) (interface{}, error)

type defaultInsert struct{}

func (d *defaultInsert) PreStart(ctx context.Context, b *testing.B, c interface{}, dataset assets.Dataset) (interface{}, error) {
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

func (d *defaultInsert) float32(ctx context.Context, b *testing.B, dataset assets.Dataset, core core.Core32) (interface{}, error) {
	train := dataset.Train()

	ids := make([]uint, 0, len(train)*10)
	for i := 0; i < 10; i++ {
		for _, vec := range train {
			id, err := core.Insert(vec)
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)
		}
	}

	return ids, nil
}

func (d *defaultInsert) float64(ctx context.Context, b *testing.B, dataset assets.Dataset, core core.Core64) (interface{}, error) {
	train := dataset.TrainAsFloat64()

	ids := make([]uint, 0, len(train)*10)
	for i := 0; i < 10; i++ {
		for _, vec := range train {
			id, err := core.Insert(vec)
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)
		}
	}

	return ids, nil
}

type defaultCreateIndex struct{}

func (d *defaultCreateIndex) PreStart(ctx context.Context, b *testing.B, c interface{}, dataset assets.Dataset) (interface{}, error) {
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

func (d *defaultCreateIndex) float32(ctx context.Context, b *testing.B, dataset assets.Dataset, core core.Core32) (interface{}, error) {
	if err := core.CreateAndSaveIndex(10000); err != nil {
		return nil, err
	}
	return []uint{}, nil
}

func (d *defaultCreateIndex) float64(ctx context.Context, b *testing.B, dataset assets.Dataset, core core.Core64) (interface{}, error) {
	if err := core.CreateAndSaveIndex(10000); err != nil {
		return nil, err
	}
	return []uint{}, nil
}
