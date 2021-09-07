package operation

import (
	"context"
	"strconv"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

func (o *operation) Insert(b *testing.B, ctx context.Context, ds assets.Dataset) int {
	b.Log("insert operation started")
	inserted := 0
	req := &payload.Insert_Request{
		Vector: &payload.Object_Vector{},
	}
	b.ResetTimer()
	b.Run("Insert", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			idx := i % ds.TrainSize()
			v, err := ds.Train(idx)
			if err != nil {
				b.Error(err)
				return
			}

			req.Vector.Id, req.Vector.Vector = strconv.FormatInt(int64(idx), 10), v.([]float32)
			_, err = o.client.Insert(ctx, req)
			if err != nil {
				st, ok := status.FromError(err)
				if !ok || st.Code() != codes.AlreadyExists {
					b.Error(err)
					return
				}
				continue
			}
			inserted++
		}
	})
	return inserted
}
