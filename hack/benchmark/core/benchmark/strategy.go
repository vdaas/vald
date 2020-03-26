package benchmark

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

// Strategy is an interface for benchmark.
type Strategy interface {
	Run(context.Context, *testing.B, assets.Dataset, []uint)
	Init(context.Context, *testing.B, assets.Dataset) error
	PreProp(context.Context, *testing.B, assets.Dataset) ([]uint, error)
	core.Closer
}
