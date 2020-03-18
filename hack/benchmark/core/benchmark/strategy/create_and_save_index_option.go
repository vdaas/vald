package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type CreateAndSaveIndexOption func(*createAndSaveIndex)

var (
	defaultCreateAndSaveIndexOptions = []CreateAndSaveIndexOption{
		WithCreateAndSaveIndexPoolSize(10000),
		WithCreateAndSaveIndexPreStart(
			func(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) {
				_, errs := ngt.BulkInsert(dataset.Train())
				for _, err := range errs {
					if err != nil {
						b.Error(err)
					}
				}
			},
		),
	}
)

func WithCreateAndSaveIndexPoolSize(size int) CreateAndSaveIndexOption {
	return func(c *createAndSaveIndex) {
		if size > 0 {
			c.poolSize = uint32(size)
		}
	}
}

func WithCreateAndSaveIndexPreStart(
	fn func(context.Context, *testing.B, ngt.NGT, assets.Dataset),
) CreateAndSaveIndexOption {
	return func(c *createAndSaveIndex) {
		if fn != nil {
			c.preStart = fn
		}
	}
}
