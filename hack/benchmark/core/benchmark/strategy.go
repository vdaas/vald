package benchmark

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/client"
)

// Strategy is an interface for benchmark.
type Strategy interface {
	Run(context.Context, *testing.B, client.Client, assets.Dataset)
}
