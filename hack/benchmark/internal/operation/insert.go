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
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

func (o *operation) Insert(b *testing.B, ds assets.Dataset) (insertedNum int) {
	b.ResetTimer()
	b.Run("Insert", func(b *testing.B) {
		ctx := b.Context()
		req := &payload.Insert_Request{
			Vector: &payload.Object_Vector{},
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: false,
			},
		}
		b.ResetTimer()

		i := 0
		for b.Loop() {
			v, err := ds.Train(i % ds.TrainSize())
			if err != nil {
				b.Fatal(err)
			}

			req.Vector.Id, req.Vector.Vector = strconv.Itoa(i), v.([]float32)

			loc, err := o.client.Insert(ctx, req)
			if err != nil {
				st, _ := status.FromError(err)
				if st.Code() != codes.AlreadyExists {
					statusError(b, int32(st.Code()), st.Message(), st.Details()...)
				}
				i++
				continue
			}
			if loc == nil {
				b.Error("returned loc is nil")
			}
			insertedNum++
			i++
		}
	})
	return insertedNum
}

func (o *operation) StreamInsert(b *testing.B, ds assets.Dataset) int {
	var insertedNum int64
	b.ResetTimer()
	b.Run("StreamInsert", func(b *testing.B) {
		ctx := b.Context()
		sc, err := o.client.StreamInsert(ctx)
		if err != nil {
			b.Fatal(err)
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		req := &payload.Insert_Request{
			Vector: &payload.Object_Vector{},
			Config: &payload.Insert_Config{
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
					// When the StreamInsert handler on the Server side returns an error, the error will be returned to Recv method.
					// In the case of multiple executions, such as benchmarking, an error will occur even if AlreadyExist occurs for some of them.
					// To prevent this, we close the stream early when an error occurs.
					return err
				}

				loc := res.GetLocation()
				if loc == nil {
					st := res.GetStatus()
					if st != nil {
						if st.GetCode() != int32(codes.AlreadyExists) {
							statusError(b, st.GetCode(), st.GetMessage(), st.GetDetails())
						}
						continue
					}
					b.Error("returned loc is nil")
					continue
				}
				atomic.AddInt64(&insertedNum, 1)
			}
		})

		i := 0
		for b.Loop() {
			v, err := ds.Train(i % ds.TrainSize())
			if err != nil {
				b.Fatal(err)
			}

			req.Vector.Id, req.Vector.Vector = strconv.Itoa(i), v.([]float32)
			err = sc.Send(req)
			if err != nil {
				b.Error(err)
			}
			i++
		}

		if err := sc.CloseSend(); err != nil {
			b.Fatal(err)
		}

		wg.Wait()
	})

	return int(insertedNum)
}
