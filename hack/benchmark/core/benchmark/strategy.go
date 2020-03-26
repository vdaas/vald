package benchmark

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
)

// Strategy is an interface for benchmark.
type Strategy interface {
	PreProcess(ctx context.Context, b *testing.B, dataset assets.Dataset) (interface{}, error)
	Run(ctx context.Context, b *testing.B, dataset assets.Dataset)
}
