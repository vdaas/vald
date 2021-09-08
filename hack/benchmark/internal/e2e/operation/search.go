package operation

import (
	"context"
	"io"
	"strconv"
	"sync"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
)

func (o *operation) Search(b *testing.B, ctx context.Context, ds assets.Dataset) {
	b.ResetTimer()
	b.Run("Search", func(b *testing.B) {
		cfg := &payload.Search_Config{
			Num:     10,
			Radius:  -1,
			Epsilon: 0.1,
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			v, err := ds.Train(i % ds.TrainSize())
			if err != nil {
				b.Error(err)
				continue
			}
			_, err = o.client.Search(ctx, &payload.Search_Request{
				Vector: v.([]float32),
				Config: cfg,
			})
			if err != nil {
				grpcError(b, err)
			}
		}
	})
}

func (o *operation) StreamSearch(b *testing.B, ctx context.Context, ds assets.Dataset) {
	b.ResetTimer()
	b.Run("StreamSearch", func(b *testing.B) {
		sc, err := o.client.StreamSearch(ctx)
		if err != nil {
			b.Fatal(err)
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		cfg := &payload.Search_Config{
			Num:     10,
			Radius:  -1,
			Epsilon: 0.1,
		}
		b.ResetTimer()

		go func() {
			defer wg.Done()

			for {
				res, err := sc.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					grpcError(b, err)
					continue
				}
				if res.GetResponse() == nil {
					b.Error("returned response is nil")
				}
			}
		}()

		for i := 0; i < b.N; i++ {
			idx := i % ds.TrainSize()
			v, err := ds.Train(idx)
			if err != nil {
				b.Error(err)
				continue
			}
			err = sc.Send(&payload.Search_Request{
				Vector: v.([]float32),
				Config: cfg,
			})
			if err != nil {
				b.Error(err)
			}
		}

		sc.CloseSend()
		wg.Wait()
	})
}

func (o *operation) SearchByID(b *testing.B, ctx context.Context, maxIdNum int) {
	b.Run("SearchByID", func(b *testing.B) {
		cfg := &payload.Search_Config{
			Num:     10,
			Radius:  -1,
			Epsilon: 0.1,
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := o.client.SearchByID(ctx, &payload.Search_IDRequest{
				Id:     strconv.Itoa(i % maxIdNum),
				Config: cfg,
			})
			if err != nil {
				grpcError(b, err)
			}
		}
	})
}

func (o *operation) StreamSearchByID(b *testing.B, ctx context.Context, maxIdNum int) {
	b.ResetTimer()
	b.Run("StreamSearchByID", func(b *testing.B) {
		sc, err := o.client.StreamSearchByID(ctx)
		if err != nil {
			b.Fatal(err)
		}
		wg := sync.WaitGroup{}
		wg.Add(1)

		cfg := &payload.Search_Config{
			Num:     10,
			Radius:  -1,
			Epsilon: 0.1,
		}
		b.ResetTimer()

		go func() {
			defer wg.Done()

			for {
				res, err := sc.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					grpcError(b, err)
					continue
				}
				if res.GetResponse() == nil {
					b.Error("returned response is nil")
				}
			}
		}()

		for i := 0; i < b.N; i++ {
			err = sc.Send(&payload.Search_IDRequest{
				Id:     strconv.Itoa(i % maxIdNum),
				Config: cfg,
			})
			if err != nil {
				b.Error(err)
			}
		}

		sc.CloseSend()
		wg.Wait()
	})
}
