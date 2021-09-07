package operation

import (
	"context"
	"strconv"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
)

func (o *operation) Search(b *testing.B, ctx context.Context, ds assets.Dataset) {
	b.Log("searchByID operation started")

	cfg := &payload.Search_Config{
		Num:     10,
		Radius:  -1,
		Epsilon: 0.1,
	}
	b.ResetTimer()
	b.Run("Search", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v, err := ds.Train(i % ds.TrainSize())
			if err != nil {
				b.Error(err)
				return
			}
			_, err = o.client.Search(ctx, &payload.Search_Request{
				Vector: v.([]float32),
				Config: cfg,
			})
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func (o *operation) SearchByID(b *testing.B, ctx context.Context, ds assets.Dataset) {
	b.Log("searchByID operation started")

	cfg := &payload.Search_Config{
		Num:     10,
		Radius:  -1,
		Epsilon: 0.1,
	}
	b.ResetTimer()
	b.Run("SearchByID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := o.client.SearchByID(ctx, &payload.Search_IDRequest{
				Id:     strconv.FormatInt(int64(i), 10),
				Config: cfg,
			})
			if err != nil {
				b.Error(err)
			}
		}
	})
}
