// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
	"errors"
	"strconv"
	"sync"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

func (o *operation) Remove(b *testing.B, maxIdNum int) {
	b.ResetTimer()
	b.Run("Remove", func(b *testing.B) {
		ctx := b.Context()
		req := &payload.Remove_Request{
			Id: &payload.Object_ID{},
			Config: &payload.Remove_Config{
				SkipStrictExistCheck: false,
			},
		}
		b.ResetTimer()

		i := 0
		for b.Loop() {
			req.Id.Id = strconv.Itoa(i % maxIdNum)
			loc, err := o.client.Remove(ctx, req)
			if err != nil {
				grpcError(b, err)
				i++
				continue
			}
			if loc == nil {
				b.Error("returned loc is nil")
			}
			i++
		}
	})
}

func (o *operation) StreamRemove(b *testing.B, maxIdNum int) {
	b.ResetTimer()
	b.Run("StreamRemove", func(b *testing.B) {
		ctx := b.Context()
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

		errgroup.Go(func() error {
			defer wg.Done()

			for {
				res, err := sc.Recv()
				if errors.Is(err, io.EOF) {
					return nil
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
		})

		i := 0
		for b.Loop() {
			req.Id.Id = strconv.Itoa(i % maxIdNum)
			err := sc.Send(req)
			if err != nil {
				b.Fatal(err)
			}
			i++
		}

		if err := sc.CloseSend(); err != nil {
			b.Fatal(err)
		}

		wg.Wait()
	})
}
