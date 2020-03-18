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
			(new(preStart)).Func,
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
