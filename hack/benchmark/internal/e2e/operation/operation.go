package operation

import (
	"context"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/client/v1/client"
	"github.com/vdaas/vald/internal/errors"
)

type Operation interface {
	Search(b *testing.B, ctx context.Context, ds assets.Dataset)
	SearchByID(b *testing.B, ctx context.Context, ds assets.Dataset)
	Insert(b *testing.B, ctx context.Context, ds assets.Dataset) int
	Remove(b *testing.B, ctx context.Context, ds assets.Dataset, maxIdNum int)
	CreateIndex(b *testing.B, ctx context.Context)
}

type operation struct {
	client   client.Client
	indexerC client.Indexer
}

func New(opts ...Option) Operation {
	o := &operation{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *operation) CreateIndex(b *testing.B, ctx context.Context) {
	b.Log("createIndex operation started")

	req := &payload.Control_CreateIndexRequest{
		PoolSize: 10000,
	}
	b.ResetTimer()
	b.Run("CreateIndex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := o.indexerC.CreateIndex(ctx, req)
			if err != nil && !errors.Is(err, errors.ErrUncommittedIndexNotFound) {
				b.Error(err)
			}
		}
	})
}
