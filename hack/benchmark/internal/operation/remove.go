// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package operation

import (
	"context"
	"strconv"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/sync"
)

func (o *operation) Remove(ctx context.Context, b *testing.B, maxIdNum int) {
	b.ResetTimer()
	b.Run("Remove", func(b *testing.B) {
		req := &payload.Remove_Request{
			Id: &payload.Object_ID{},
			Config: &payload.Remove_Config{
				SkipStrictExistCheck: false,
			},
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			req.Id.Id = strconv.Itoa(i % maxIdNum)
			loc, err := o.client.Remove(ctx, req)
			if err != nil {
				grpcError(b, err)
				continue
			}
			if loc == nil {
				b.Error("returned loc is nil")
			}
		}
	})
}

func (o *operation) StreamRemove(ctx context.Context, b *testing.B, maxIdNum int) {
	b.ResetTimer()
	b.Run("StreamRemove", func(b *testing.B) {
		sc, err := o.client.StreamRemove(ctx)
		if err != nil {
			b.Fatal(err)
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		req := &payload.Remove_Request{
			Id: &payload.Object_ID{},
			Config: &payload.Remove_Config{
				SkipStrictExistCheck: false,
			},
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

				loc := res.GetLocation()
				if loc == nil {
					st := res.GetStatus()
					if st != nil {
						statusError(b, st.GetCode(), st.GetMessage(), st.GetDetails())
						continue
					}

					b.Error("returned response is nil")
					continue
				}
			}
		}()

		for i := 0; i < b.N; i++ {
			req.Id.Id = strconv.Itoa(i % maxIdNum)
			err := sc.Send(req)
			if err != nil {
				b.Fatal(err)
			}
		}

		if err := sc.CloseSend(); err != nil {
			b.Fatal(err)
		}
		wg.Wait()
	})
}
