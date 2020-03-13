// Package strategy provides strategy for e2e testing functions
package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/internal/client"
)

type createIndex struct {
	poolSize uint32
	client.Indexer
}

func NewCreateIndex(opts ...CreateIndexOption) e2e.Strategy {
	ci := new(createIndex)
	for _, opt := range append(defaultCreateIndexOptions, opts...) {
		opt(ci)
	}
	return ci
}

func (ci *createIndex) Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	b.Run("CreateIndex", func(bb *testing.B) {
		for i := 0; i < bb.N; i++ {
			ci.do(ctx, b)
		}
	})
}

func (ci *createIndex) do(ctx context.Context, b *testing.B) {
	if err := ci.Indexer.CreateIndex(ctx, &client.ControlCreateIndexRequest{
		PoolSize: ci.poolSize,
	}); err != nil {
		b.Error(err)
	}
}
