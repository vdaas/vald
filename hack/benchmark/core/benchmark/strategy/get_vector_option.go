package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type GetVectorOption func(*getVector)

var (
	defaultGetVectorOptions = []GetVectorOption{
		WithGetVectorPreStart(
			func(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) (ids []uint) {
				ids = (new(preStart)).Func(ctx, b, ngt, dataset)
				if err := ngt.CreateIndex(10000); err != nil {
					b.Error(err)
				}
				return ids
			},
		),
	}
)

func WithGetVectorPreStart(fn PreStart) GetVectorOption {
	return func(g *getVector) {
		if fn != nil {
			g.preStart = fn
		}
	}
}
