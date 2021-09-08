package operation

import (
	"context"
	"io"
	"strconv"
	"sync"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

func (o *operation) Insert(b *testing.B, ctx context.Context, ds assets.Dataset) (insertedNum int) {
	b.ResetTimer()
	b.Run("Insert", func(b *testing.B) {
		req := &payload.Insert_Request{
			Vector: &payload.Object_Vector{},
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			idx := i % ds.TrainSize()
			v, err := ds.Train(idx)
			if err != nil {
				b.Error(err)
				continue
			}

			req.Vector.Id, req.Vector.Vector = strconv.Itoa(idx), v.([]float32)

			loc, err := o.client.Insert(ctx, req)
			if err != nil {
				st, ok := status.FromError(err)
				if !ok || st.Code() != codes.AlreadyExists {
					statusError(b, st)
				}
				continue
			}
			if loc == nil {
				b.Error("returned loc is nil")
			}
			insertedNum++
		}
	})
	return insertedNum
}

func (o *operation) StreamInsert(b *testing.B, ctx context.Context, ds assets.Dataset) (insertedNum int) {
	b.ResetTimer()
	b.Run("StreamInsert", func(b *testing.B) {
		sc, err := o.client.StreamInsert(ctx)
		if err != nil {
			b.Fatal(err)
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		req := &payload.Insert_Request{
			Vector: &payload.Object_Vector{},
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
					st, ok := status.FromError(err)
					if !ok || st.Code() != codes.AlreadyExists {
						statusError(b, st)
					}
					continue
				}

				loc := res.GetLocation()
				if loc == nil {
					b.Error("returned loc is nil")
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

			req.Vector.Id, req.Vector.Vector = strconv.Itoa(idx), v.([]float32)
			err = sc.Send(req)
			if err != nil {
				b.Error(err)
				continue
			}
			insertedNum++
		}

		sc.CloseSend()
		wg.Wait()
	})

	return insertedNum
}
