package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
)

type RemoveOption func(*remove)

var (
	defaultRemoveOptions = []RemoveOption{
		WithRemovePreStart(
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

func WithRemovePreStart(fn PreStart) RemoveOption {
	return func(r *remove) {
		if fn != nil {
			r.preStart = fn
		}
	}
}
