package operation

import (
	"context"
	"strconv"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
)

func (o *operation) Remove(b *testing.B, ctx context.Context, ds assets.Dataset, maxIdNum int) {
	b.Log("remove operation started")

	req := &payload.Remove_Request{
		Id: &payload.Object_ID{},
		Config: &payload.Remove_Config{
			SkipStrictExistCheck: false,
		},
	}
	b.ResetTimer()
	b.Run("Remove", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if i < maxIdNum {
				req.Id.Id = strconv.FormatInt(int64(i), 10)
				_, err := o.client.Remove(ctx, req)
				if err != nil {
					b.Error(err)
				}
			}
		}
	})
}
