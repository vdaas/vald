package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/helper/e2e"
	"github.com/vdaas/vald/internal/client"
)

type streamInsert struct{}

func NewStreamInsert(opts ...StreamInsertOption) e2e.Strategy {
	s := new(streamInsert)

	for _, opt := range append(defaultStreamInsertOptions, opts...) {
		opt(s)
	}

	return s
}

func (sisrt *streamInsert) dataProvider() func() *client.ObjectVector {
	return nil
}

func (sisrt *streamInsert) Run(ctx context.Context, b *testing.B, client client.Client, dataset assets.Dataset) {
}
