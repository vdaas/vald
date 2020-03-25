package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
)

type GetVectorOption func(*getVector)

var (
	defaultGetVectorOptions = []GetVectorOption{
		WithGetVectorPreStart(
			func(ctx context.Context, b *testing.B, c interface{}, dataset assets.Dataset) (interface{}, error) {
				ids, err := (new(defaultInsert)).PreStart(ctx, b, c, dataset)
				if err != nil {
					return nil, err
				}

				_, err = (new(defaultCreateIndex)).PreStart(ctx, b, c, dataset)
				if err != nil {
					return nil, err
				}

				return ids, nil
			},
		),
	}
)

func WithGetVectorPreStart(fn PreStart) GetVectorOption {
	return func(g *getVector) {
		if g.preStart != nil {
			g.preStart = fn
		}
	}
}
