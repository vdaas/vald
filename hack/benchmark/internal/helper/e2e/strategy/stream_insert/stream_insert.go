package streaminsert

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/helper/e2e"
	"github.com/vdaas/vald/internal/client"
)

type streaminsert struct{}

func New() e2e.Strategy {
	return nil
}

func (sisrt *streaminsert) dataProvider() func() *client.ObjectVector {
	return nil
}

func (sisrt *streaminsert) Run(ctx context.Context, b *testing.B, dataset assets.Dataset) error {
	return nil
}
