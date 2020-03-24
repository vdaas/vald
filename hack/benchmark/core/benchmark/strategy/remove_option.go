package stratedy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
)

type RemoveOption func(*remove)

var (
	defaultOptions = []RemoveOption{}
)

func WithPreStart(
	fn func(context.Context, *testing.B, assets.Dataset) (interface{}, error),
) RemoveOption {
	return func(d *remove) {
		if fn != nil {
			d.preStart = fn
		}
	}
}
