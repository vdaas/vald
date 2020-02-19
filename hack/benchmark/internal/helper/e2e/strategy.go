package e2e

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/client"
)

type Strategy interface {
	Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) error
}
