package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type BulkRemoveOption func(*bulkRemove)

var (
	defaultBulkRemoveOptions = []BulkRemoveOption{
		WithBulkRemoveChunkSize(1000),
		WithBulkRemovePrestart(
			func(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) (ids []uint) {
				ids = make([]uint, 0, len(dataset.Train()))
				isrtIds, errs := ngt.BulkInsert(dataset.Train())
				for i, err := range errs {
					if err != nil {
						b.Error(err)
					} else {
						ids = append(ids, isrtIds[i])
					}
				}
				return
			},
		),
	}
)

func WithBulkRemoveChunkSize(chunk int) BulkRemoveOption {
	return func(br *bulkRemove) {
		if chunk > 0 {
			br.chunkSize = chunk
		}
	}
}

func WithBulkRemovePrestart(
	fn func(context.Context, *testing.B, ngt.NGT, assets.Dataset) []uint,
) BulkRemoveOption {
	return func(br *bulkRemove) {
		if fn != nil {
			br.fn = fn
		}
	}
}
