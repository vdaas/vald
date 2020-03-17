package benchmark

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

// Strategy is an interface for benchmark.
type Strategy interface {
	Run(context.Context, *testing.B, ngt.NGT, assets.Dataset)
}
